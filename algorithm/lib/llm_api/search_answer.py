from AskOnce.algorithm.lib.llm_api.llm_base_api import LLMBaseAPI
from AskOnce.algorithm.lib.llm_api.generate_outlines import GenerateOutlines
from AskOnce.algorithm.lib.llm_api.question_process import QuestionProcess
class SearchAnswer (LLMBaseAPI):
    
    def simplify_answer(self,question,search_result,stream=False):
        print('对问题进行简答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题进行回答，要求：
1. 回答内容简洁全面，只需要回答问题关注的信息，不要做过多的解释。

问题：{question}'''
            prompt_ok = prompt.format(question=question)
        else:
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的相关内容），对用户输入的问题进行回答，要求：
1. 回答内容简洁全面，只需要回答问题关注的信息，不要做过多的解释。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。

相关内容：{answers}

问题：{question}'''
            answers = []
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.1,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def answer_with_history_messages(self,question,search_result,history_messages,stream=False):
        print('追问',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result,max_search_content_length=3000,max_chars=3000)
        if len(all_search_result)==0:
            prompt = '''
'''
        else:
            prompt = '''
'''
        
        prompt_ok  = ''
        history_messages.append({'role':'user','content':prompt_ok})
        return_result = {
            'prompt' :history_messages,
            'temperature':0.3,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def detailed_answer(self,question,search_result,stream=False):
        print('对问题进行段落划分的复杂回答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多，只需要回答问题关注的信息，不要做过多的解释。
2. 使用markdown格式，回答部分进行章节划分，不要用一级标题，每个段落的文字要内容丰富，字数多。

问题：{question}'''
            prompt_ok = prompt.format(question=question)
        else:
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多，只需要回答问题关注的信息，不要做过多的解释。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
4. 使用markdown格式，回答部分进行章节划分，不要用一级标题，每个段落的文字要内容丰富，字数多。

相关内容：{answers}

问题：{question}'''
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.3,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    
    def detailed_no_chapter_answer(self,question,search_result,stream=False):
        print('对问题进行无段落划分的复杂回答:',question)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        answers = []
        if len(all_search_result)==0:
            print('搜索结果都过滤掉了')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多。
2. 使用markdown格式，但不要有章节层级，使用大的段落。
3. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。

问题：{question}''' 
            prompt_ok = prompt.format(question=question)
        else:
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题进行回答，要求：
1. 回答内容丰富详细，字数多。
2. 只使用相关内容的信息，使用第三人称的表达方法。
3. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
4. 使用markdown格式，但不要有章节层级，使用大的段落。
5. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。

相关内容：{answers}

问题：{question}''' 
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt_ok = prompt.format(question=question,answers='\n'.join(answers))
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.3,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
    # 专业的回答
    # 1. 会生成三个相关的问题，并进行回答
    # 2. 专业的回答会先生成目录大纲，逐个目录大纲进行内容丰富。
    
    
    def professional_answer_one_chapter(self,question,chapter,directory_str,search_result,stream=False):
        print('专家回答章节:',chapter)
        print('搜索结果长度',len(search_result))
        all_search_result = self.select_search_result(search_result)
        
        if len(all_search_result)==0:
            print('无搜索结果')
            prompt = '''你现在是一个AI智能助手，对用户输入的问题，按照大纲，回答大纲中的指定章节，要求：
1. 回答内容丰富详细，字数多。
2. 使用markdown格式，但不要有章节层级，使用大的段落。
3. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。
4. 回答的内容只限于指定章节的内容，不要与其他章节内容冲突。

问题：{question}

大纲：
{directory_str}

待回答的指定章节：{chapter}

回答内容：
'''   
            prompt_ok = prompt.format(question=question,directory_str=directory_str,chapter=chapter)
        else:
            answers=[]
            for item in all_search_result:
                answers.append('标题：'+item['title']+'\n部分内容：'+item['answer_for_question'].strip())
            prompt = '''你现在是一个AI智能搜索引擎，根据相关内容（根据问题搜索出来的结果），对用户输入的问题，按照大纲，回答大纲中的指定章节，要求：
1. 回答内容丰富详细，字数多。
2. 使用markdown格式，但不要有章节层级，使用大的段落。
3. 只需要回答问题关注的信息，不要做过多的解释和内容的延伸。
4. 回答的内容只限于指定章节的内容，不要与其他章节内容冲突。
5. 不要输出参考资料，不要输出相关信息的来源，例如标题等，也不要输出“根据已知信息”相关的表述。
6. 只使用相关内容的信息，使用第三人称的表达方法。

相关内容：{search_result_str}

问题：{question}

大纲：
{directory_str}

待回答的指定章节：{chapter}

回答内容：
'''         
            prompt_ok = prompt.format(question=question,directory_str=directory_str,chapter=chapter,search_result_str='\n'.join(answers))
        
        return_result = {
            'prompt' :prompt_ok,
            'temperature':0.4,
            'max_tokens':2048
        }
        if stream:
            return self.ask_llm_stream(**return_result)
        else:
            return self.ask_llm(**return_result)
            
    def professional_answer(self,question,search_result,stream=False):
        print('对问题进行专业回答:',question)
        print('搜索结果长度',len(search_result))
        # 先生成回答的大纲
        outlines_model  =  GenerateOutlines(platform_api_url = self.platform_api_url,
                                            api_key  =self.api_key,
                                            model_name = self.model_name)
        directory = outlines_model.generate_outlines_by_question_and_search_result(question,search_result)
        all_result = ''
        if len(directory) ==0 :
            detailed_answer_result = self.detailed_answer(question,search_result,stream)
            if stream:
                for item in detailed_answer_result:
                    all_result+=item
                    yield item
                yield '\n'
                all_result+='\n'
            else:
                all_result+=detailed_answer_result+'\n'
                yield detailed_answer_result

        directory_str = '\n'.join([ item['title_level']+'# '+item['content'] for item in directory ])
        print('大纲：\n',directory_str)
        # 根据大纲逐步生成回答
        for index,directory_item in enumerate(directory):
            
            if index < len(directory)-1:
                all_result += directory_item['title_level']+'# '+directory_item['content']+'\n'
                yield directory_item['title_level']+'# '+directory_item['content']+'\n'
                if len(directory[index]['title_level']) < len(directory[index+1]['title_level']):
                    continue
            else:
                all_result += directory_item['title_level']+'# '+directory_item['content']+'\n'
                yield directory_item['title_level']+'# '+directory_item['content']+'\n'
                
                    
            one_chapter_result = self.professional_answer_one_chapter(question,directory_item['content'],directory_str,search_result,stream)
            if stream:
                for item in one_chapter_result:
                    all_result+=item
                    yield item
                yield '\n'
                all_result+='\n'
            else:
                all_result+=one_chapter_result+'\n'
                yield one_chapter_result
                yield '\n'
        
        
        # 生成额外的问题
        question_process_model = QuestionProcess(platform_api_url = self.platform_api_url,
                                            api_key  =self.api_key,
                                            model_name = self.model_name)
        more_question_or_topic_list = question_process_model.generate_more_related_question(question,all_result)
        # 对额外的问题生成回答
        print('额外的问题',more_question_or_topic_list)
        yield '\n---\n'
        
        for more_question in more_question_or_topic_list:
            yield '> ## '+more_question+'\n'
            more_question_answer = self.detailed_no_chapter_answer(more_question,search_result,stream)
            yield '> '
            if stream:
                for item in more_question_answer:
                    if item =='\n':
                        item += '\n> '
                    yield item
            else:
                yield '> '+more_question_answer+'\n'
            yield '\n'
        
    
    